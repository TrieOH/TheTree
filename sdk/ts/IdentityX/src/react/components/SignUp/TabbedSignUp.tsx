import { useState } from "react";
import { SignUp, type SignUpProps } from "./SignUp";

export interface TabbedSignUpProps extends Omit<SignUpProps, 'flow_id'> {
  flowIds: { label: string; value: string; }[];
}

export function TabbedSignUp({ flowIds, ...rest }: TabbedSignUpProps) {
  const [activeFlowId, setActiveFlowId] = useState<string>(flowIds[0]?.value || "");

  if (!flowIds || flowIds.length === 0) {
    return (
      <div className="trieoh-card trieoh-card--full-rounded">
        <p className="trieoh-tabbed-signup__empty">Nenhuma opção de cadastro disponível.</p>
      </div>
    );
  }

  return (
    <div className="trieoh-tabbed-signup">
      <div className="trieoh-tabbed-signup__header">
        {flowIds.map((flow) => (
          <button
            key={flow.value}
            className={`trieoh-tabbed-signup__tab ${activeFlowId === flow.value ? "active" : ""}`}
            onClick={() => setActiveFlowId(flow.value)}
            type="button"
            aria-selected={activeFlowId === flow.value}
            role="tab"
          >
            <span className="trieoh-tabbed-signup__tab-text">{flow.label}</span>
          </button>
        ))}
      </div>
      <hr />

      <div className="trieoh-tabbed-signup__content">
        <SignUp flow_id={activeFlowId} {...rest} />
      </div>
    </div>
  );
}