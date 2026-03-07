export const addActivityPermission = (activityId) => ({
    permission_type: "activity",
    activity_id: activityId,
    product_id: null,
    checkpoint_id: null,
});

export const addProductPermission = (productId) => ({
    permission_type: "product",
    activity_id: null,
    product_id: productId,
    checkpoint_id: null,
});

export const addCheckpointPermission = (checkpointId) => ({
    permission_type: "checkpoint",
    activity_id: null,
    product_id: null,
    checkpoint_id: checkpointId,
});