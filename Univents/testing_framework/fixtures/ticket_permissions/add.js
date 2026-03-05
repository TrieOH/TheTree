export const addActivityPermission = (activityId) => ({
    ticket_scope_id: "", // fill at test time
    permission_type: "activity",
    activity_id: activityId,
    product_id: null,
    checkpoint_id: null,
});

export const addProductPermission = (productId) => ({
    ticket_scope_id: "",
    permission_type: "product",
    activity_id: null,
    product_id: productId,
    checkpoint_id: null,
});

export const addCheckpointPermission = (checkpointId) => ({
    ticket_scope_id: "",
    permission_type: "checkpoint",
    activity_id: null,
    product_id: null,
    checkpoint_id: checkpointId,
});