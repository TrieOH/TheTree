import type { FieldI } from "#/features/fields/model";
import type { FormI } from "#/features/forms/model";
import type { StepI } from "#/features/steps/model";
import type { Answer, SubmitRequest } from "@trieoh/informd-models";
import { z } from "zod";

export const answerSchema = z.object({
  id: z.string(), // i don't need on submit
  response_id: z.string(), // i don't need on submit
  field_id: z.string().optional(),
  answer: z.any().optional(),
  answered_at: z.string(), // i don't need on submit
  updated_at: z.string().optional(), // i don't need on submit
}) satisfies z.ZodType<Answer>;

export const submitRequestSchema = z.object({
  email: z.email().optional().or(z.literal("")),
  answers: z.array(answerSchema),
}) satisfies z.ZodType<SubmitRequest>;

export type SubmitRequestI = SubmitRequest;
export type AnswerI = Answer;

export interface FullFormI {
  form: FormI;
  steps: FullStepI[];
}

export interface FullStepI {
  step: StepI;
  fields: FullFieldI[] | null;
}

export interface FullFieldI {
  field: FieldI;
  answers: FullAnswerI[];
}

export interface FullAnswerI {
  answer: AnswerI;
  responder: string;
}

export interface SubmissionSummaryI {
  responder: string;
  completed_at: string;
  step_id: string;
  answers: Record<string, string>;
}

export function deriveSubmissions(fullForm: FullFormI): SubmissionSummaryI[] {
  const byResponder = new Map<string, FullAnswerI[]>();

  for (const step of fullForm.steps) {
    for (const field of step.fields ?? []) {
      for (const fa of field.answers) {
        const list = byResponder.get(fa.responder) ?? [];
        list.push(fa);
        byResponder.set(fa.responder, list);
      }
    }
  }

  const submissions: SubmissionSummaryI[] = [];

  for (const [responder, answers] of byResponder) {
    const timestamps = answers.map((a) => new Date(a.answer.answered_at).getTime());
    const lastTimestamp = Math.max(...timestamps);
    const completedAt = new Date(lastTimestamp).toISOString();

    let lastStepId = "";
    let maxPosition = -1;

    for (const step of fullForm.steps) {
      const hasAny = (step.fields ?? []).some((f) =>
        f.answers.some(
          (a) =>
            a.responder === responder &&
            a.answer.answer !== undefined &&
            a.answer.answer !== ""
        )
      );
      if (hasAny && step.step.position_hint > maxPosition) {
        maxPosition = step.step.position_hint;
        lastStepId = step.step.id;
      }
    }

    const answerMap: Record<string, string> = {};
    for (const a of answers) {
      if (a.answer.field_id && a.answer.answer !== undefined) {
        answerMap[a.answer.field_id] = a.answer.answer;
      }
    }

    submissions.push({
      responder,
      completed_at: completedAt,
      step_id: lastStepId,
      answers: answerMap,
    });
  }

  return submissions.sort(
    (a, b) => new Date(b.completed_at).getTime() - new Date(a.completed_at).getTime()
  );
}