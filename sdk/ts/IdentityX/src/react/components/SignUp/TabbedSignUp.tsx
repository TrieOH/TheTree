import { useState } from "react";
import { SignUp, type SignUpProps } from "./SignUp";
import type { FieldDefinitionResultI } from "../../../types/fields-types";

export interface TabbedFlowI {
  label: string;
  value: string;
  fields: FieldDefinitionResultI[];
}

export interface TabbedSignUpProps extends Omit<SignUpProps, 'flow_id' | 'fields'> {
  flowIds: TabbedFlowI[];
}

export function TabbedSignUp({ flowIds, ...rest }: TabbedSignUpProps) {
  const [activeFlowId, setActiveFlowId] = useState<string>(flowIds[0]?.value || "");

  if (!flowIds || flowIds.length === 0) return null;

  const activeFlow = flowIds.find(f => f.value === activeFlowId);

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
        <SignUp 
          key={activeFlowId} 
          flow_id={activeFlowId} 
          {...rest} 
          fields={activeFlow?.fields} 
        />
      </div>
    </div>
  );
}