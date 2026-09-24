import { render, screen, fireEvent, cleanup, waitFor } from "@solidjs/testing-library";
import { describe, it, expect, vi, afterEach, beforeEach } from "vitest";
import {
  FulfillSignatureRequestView,
  RevokeSignatureView,
} from "@/features/signatures/ui";
import {
  parseJwtPayload,
  signatureRequestTokenClaimsSchema,
  signatureRevocationTokenClaimsSchema,
  type SignatureRequestTokenClaims,
  type SignatureRevocationTokenClaims,
} from "@/features/signatures/model";

// Mock TanStack Solid Router Link
vi.mock("@tanstack/solid-router", () => ({
  Link: (props: { children?: unknown; class?: string; to?: string; [key: string]: unknown }) => {
    return (
      <a href={props.to ?? "#"} class={props.class}>
        {props.children as any}
      </a>
    );
  },
}));

let mockQueryReturn = {
  data: null as any,
  isLoading: false,
  isError: false,
};

// Mock Query Client / Front Core
vi.mock("@trieoh/front-core-solid", () => ({
  useQuery: vi.fn(() => () => mockQueryReturn),
  useQueryClient: vi.fn(() => ({
    invalidateQueries: vi.fn(),
  })),
}));

// Mock Mutations
vi.mock("@/features/signatures/api/mutations", () => ({
  useFulfillSignatureRequestMutation: vi.fn(() => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
  })),
  useDenySignatureRequestMutation: vi.fn(() => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
  })),
  useRevokeSignatureMutation: vi.fn(() => ({
    mutateAsync: vi.fn().mockResolvedValue({}),
  })),
}));

describe("Signatures Fulfill & Revoke Views", () => {
  beforeEach(() => {
    mockQueryReturn = {
      data: null,
      isLoading: false,
      isError: false,
    };
  });

  afterEach(() => {
    cleanup();
  });

  describe("parseJwtPayload with Zod validation", () => {
    it("safely decodes and validates a valid JWT payload for requests", () => {
      const payloadObj = {
        request_id: "req-123",
        edition_id: "ed-456",
        exp: 1900000000,
      };
      const base64UrlPayload = btoa(JSON.stringify(payloadObj))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "");
      const fakeJwt = `eyJhbGciOiJIUzI1NiJ9.${base64UrlPayload}.mocksignature`;

      const claims = parseJwtPayload<SignatureRequestTokenClaims>(
        fakeJwt,
        signatureRequestTokenClaimsSchema,
      );
      expect(claims).not.toBeNull();
      expect(claims?.request_id).toBe("req-123");
      expect(claims?.edition_id).toBe("ed-456");
      expect(claims?.exp).toBe(1900000000);
    });

    it("safely decodes and validates a valid JWT payload for revocations", () => {
      const payloadObj = {
        signature_id: "sig-789",
        edition_id: "ed-456",
        exp: 1900000000,
      };
      const base64UrlPayload = btoa(JSON.stringify(payloadObj))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "");
      const fakeJwt = `eyJhbGciOiJIUzI1NiJ9.${base64UrlPayload}.mocksignature`;

      const claims = parseJwtPayload<SignatureRevocationTokenClaims>(
        fakeJwt,
        signatureRevocationTokenClaimsSchema,
      );
      expect(claims).not.toBeNull();
      expect(claims?.signature_id).toBe("sig-789");
      expect(claims?.edition_id).toBe("ed-456");
      expect(claims?.exp).toBe(1900000000);
    });

    it("returns null for invalid or malformed tokens or schemas that fail validation", () => {
      expect(parseJwtPayload("")).toBeNull();
      expect(parseJwtPayload("invalidtoken")).toBeNull();
      expect(parseJwtPayload("invalid.token")).toBeNull();

      // Missing required request_id
      const invalidObj = { edition_id: "ed-456" };
      const base64 = btoa(JSON.stringify(invalidObj));
      const token = `header.${base64}.sig`;
      expect(
        parseJwtPayload(token, signatureRequestTokenClaimsSchema),
      ).toBeNull();
    });
  });

  describe("FulfillSignatureRequestView", () => {
    it("renders invalid token state when no token is provided", () => {
      render(() => <FulfillSignatureRequestView token="" />);

      expect(screen.getByText("Link Inválido ou Incompleto")).toBeInTheDocument();
      expect(
        screen.getByText(/Não foi possível identificar o token de autenticação/),
      ).toBeInTheDocument();
    });

    it("renders page layout when valid unexpired token is provided", () => {
      const payloadObj = {
        request_id: "req-123",
        edition_id: "ed-456",
        exp: Math.floor(Date.now() / 1000) + 86400, // tomorrow
      };
      const base64UrlPayload = btoa(JSON.stringify(payloadObj))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "");
      const fakeJwt = `eyJhbGciOiJIUzI1NiJ9.${base64UrlPayload}.mocksignature`;

      mockQueryReturn = {
        data: {
          id: "req-123",
          edition_id: "ed-456",
          signatory_name: "Maria da Silva",
          signatory_title: "Diretora Acadêmica",
          status: "pending",
          expires_at: new Date(Date.now() + 86400000).toISOString(),
        },
        isLoading: false,
        isError: false,
      };

      render(() => <FulfillSignatureRequestView token={fakeJwt} />);

      expect(screen.getByText("Assinatura Oficial de Certificados")).toBeInTheDocument();
      expect(screen.getByText("Desenhar")).toBeInTheDocument();
      expect(screen.getByText("Upload de arquivo")).toBeInTheDocument();
      expect(screen.getByText("Confirmar e Assinar")).toBeInTheDocument();
    });

    it("toggles between draw and upload mode with drag & drop support", async () => {
      const payloadObj = {
        request_id: "req-123",
        edition_id: "ed-456",
        exp: Math.floor(Date.now() / 1000) + 86400,
      };
      const base64UrlPayload = btoa(JSON.stringify(payloadObj))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "");
      const fakeJwt = `eyJhbGciOiJIUzI1NiJ9.${base64UrlPayload}.mocksignature`;

      mockQueryReturn = {
        data: {
          id: "req-123",
          edition_id: "ed-456",
          signatory_name: "Maria da Silva",
          status: "pending",
          expires_at: new Date(Date.now() + 86400000).toISOString(),
        },
        isLoading: false,
        isError: false,
      };

      render(() => <FulfillSignatureRequestView token={fakeJwt} />);

      const uploadTab = screen.getByText("Upload de arquivo");
      fireEvent.click(uploadTab);

      await waitFor(() => {
        expect(
          screen.getByText(/Arraste e solte o arquivo aqui, ou clique para selecionar/),
        ).toBeInTheDocument();
      });

      // Simulate file drop
      const dropzone = screen.getByText(/Arraste e solte o arquivo aqui/);
      const fakeFile = new File(["dummy signature content"], "my-signature.png", {
        type: "image/png",
      });

      fireEvent.dragEnter(dropzone);
      await waitFor(() => {
        expect(screen.getByText("Solte o arquivo de imagem aqui")).toBeInTheDocument();
      });

      fireEvent.drop(dropzone, {
        dataTransfer: {
          files: [fakeFile],
        },
      });

      await waitFor(() => {
        expect(screen.getByText("my-signature.png")).toBeInTheDocument();
        expect(screen.getByText("Trocar")).toBeInTheDocument();
        expect(screen.getByText("Remover")).toBeInTheDocument();
      });
    });
  });

  describe("RevokeSignatureView", () => {
    it("renders invalid token state when no token is provided", () => {
      render(() => <RevokeSignatureView token="" />);

      expect(screen.getByText("Link Inválido ou Incompleto")).toBeInTheDocument();
      expect(
        screen.getByText(/Não foi possível validar o token de revogação/),
      ).toBeInTheDocument();
    });

    it("renders signature preview and supports revoking with valid token", async () => {
      const payloadObj = {
        signature_id: "sig-789",
        edition_id: "ed-456",
        exp: Math.floor(Date.now() / 1000) + 86400,
      };
      const base64UrlPayload = btoa(JSON.stringify(payloadObj))
        .replace(/\+/g, "-")
        .replace(/\//g, "_")
        .replace(/=+$/, "");
      const fakeJwt = `eyJhbGciOiJIUzI1NiJ9.${base64UrlPayload}.mocksignature`;

      mockQueryReturn = {
        data: {
          id: "sig-789",
          edition_id: "ed-456",
          signatory_name: "Maria da Silva",
          signatory_title: "Diretora Acadêmica",
          signatory_email: "maria@example.com",
          image_url: "https://example.com/signature.png",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        isLoading: false,
        isError: false,
      };

      render(() => <RevokeSignatureView token={fakeJwt} />);

      expect(screen.getByText("Revogação de Assinatura Digital")).toBeInTheDocument();
      expect(screen.getByText("Maria da Silva")).toBeInTheDocument();
      expect(screen.getByText("Diretora Acadêmica")).toBeInTheDocument();
      expect(screen.getByText("Revogar Assinatura")).toBeInTheDocument();

      const checkbox = screen.getByRole("checkbox") as HTMLInputElement;
      fireEvent.click(checkbox);

      const revokeButton = screen.getByRole("button", { name: /revogar assinatura/i });
      await waitFor(() => {
        expect(revokeButton).not.toBeDisabled();
      });

      fireEvent.click(revokeButton);

      await waitFor(() => {
        expect(screen.getByText("Assinatura Revogada")).toBeInTheDocument();
        expect(screen.getByText("Fechar")).toBeInTheDocument();
      });
    });
  });
});
