const now = new Date();
const inMinutes = (offset) => new Date(now.getTime() + offset * 60 * 1000).toISOString();

export const createRustActivity = {
    edition_scope_id: "", // fill in at test time
    title: "Introduction to Rust",
    description: "A beginner friendly introduction to Rust",
    location: "Auditório Principal",
    starts_at: inMinutes(11),
    ends_at: inMinutes(14),
    presenter_name: "João Silva",
    token_cost: 0,
    has_capacity: true,
    capacity: 50,
    difficulty: "beginner",
};

export const createKubernetesActivity = {
    edition_scope_id: "", // fill in at test time
    title: "Advanced uses of Kubernetes",
    description: "A deep dive into pods and K8s",
    location: "Auditório Secundário",
    starts_at: inMinutes(11),
    ends_at: inMinutes(14),
    presenter_name: "Marcelo Diniz",
    token_cost: 1,
    has_capacity: true,
    capacity: 50,
    difficulty: "advanced",
};