import { createFileRoute, useParams } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { Button } from '#/shared/ui/shadcn/button'
import { ConfirmModal } from '#/widgets/modal/modal'
import { webhookCreateSchema } from '#/features/webhooks/model'
import { toast } from 'sonner'
import FormModal from '#/widgets/modal/form-modal'
import type { WebhookCreateI, WebhookCreateResponseI, WebhookI } from '#/features/webhooks/model'
import { useState } from 'react'
import { allWorkspaceWebhooksQueryOptions, registerWebhookOnWorkspaceFn, deleteWebhookOnWorkspaceFn } from '#/features/webhooks/api'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { WebhookList } from '#/features/webhooks/ui/webhook-list'
import { WebhookCreatedModal } from '#/features/webhooks/ui/webhook-created-modal'

export const Route = createFileRoute('/admin/$name/webhooks')({
  component: RouteComponent,
})

function RouteComponent() {
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [deleteWebhookId, setDeleteWebhookId] = useState<string | null>(null)
  const [newlyCreatedWebhook, setNewlyCreatedWebhook] = useState<WebhookCreateResponseI | null>(null)

  const { name } = useParams({ from: '/admin/$name' })
  const queryClient = useQueryClient();
  const { data: webhooks = [], isLoading } = useQuery(allWorkspaceWebhooksQueryOptions(name))

  const { mutate: createWebhook, isPending: isPendingCreate } = useMutation({
    mutationFn: (data: WebhookCreateI) => registerWebhookOnWorkspaceFn(data, name),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspaceWebhooksQueryOptions(name).queryKey,
          (old: WebhookI[] = []) => [response.data, ...old],
        )
        setIsCreateOpen(false)
        setNewlyCreatedWebhook(response.data)
        toast.success("Webhook created successfully")
      }
    },
  })

  const { mutate: deleteWebhook } = useMutation({
    mutationFn: (id: string) => deleteWebhookOnWorkspaceFn(name, id),
    onSuccess: (response) => {
      if (response.success) {
        queryClient.setQueryData(
          allWorkspaceWebhooksQueryOptions(name).queryKey,
          (old: WebhookI[] = []) => old.filter(w => w.id !== deleteWebhookId)
        );
        setDeleteWebhookId(null)
        toast.success(response.message)
      }
    },
  })

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4">
        <div className="space-y-1">
          <h2 className="text-2xl md:text-3xl font-black uppercase tracking-tighter">Webhooks</h2>
          <p className="text-muted-foreground text-sm uppercase tracking-wider font-bold opacity-70">Receive real-time notifications about payment events.</p>
        </div>

        <Button
          onClick={() => setIsCreateOpen(true)}
          className="rounded-none gap-2 h-10 font-black uppercase tracking-widest transition-all"
        >
          <Plus className="w-4 h-4" />
          Add Endpoint
        </Button>
      </div>

      <WebhookList
        webhooks={webhooks}
        isLoading={isLoading}
        onDelete={setDeleteWebhookId}
      />

      {/* Create Modal */}
      <FormModal<WebhookCreateI>
        title="Add Webhook Endpoint"
        description="Enter the URL where you want to receive payment events."
        buttonTitle="Add Endpoint"
        schema={webhookCreateSchema}
        formId="create-webhook-form"
        isOpen={isCreateOpen}
        onClose={() => setIsCreateOpen(false)}
        onSubmit={createWebhook}
        fields={[
          {
            name: "url",
            label: "Endpoint URL",
            type: "text",
            placeholder: "https://api.your-app.com/webhooks"
          }
        ]}
        disabled={isPendingCreate}
      />

      {/* New Webhook Result Modal */}
      <WebhookCreatedModal
        webhook={newlyCreatedWebhook}
        isOpen={!!newlyCreatedWebhook}
        onClose={() => setNewlyCreatedWebhook(null)}
      />

      {/* Delete Confirmation Modal */}
      {deleteWebhookId && <ConfirmModal
        isOpen={!!deleteWebhookId}
        onClose={() => setDeleteWebhookId(null)}
        onConfirm={() => deleteWebhook(deleteWebhookId)}
        title="Delete Webhook"
        description="Are you sure you want to delete this webhook endpoint? You will stop receiving events at this URL immediately."
        confirmText="Delete Endpoint"
        variant="destructive"
      />}
    </div>
  )
}
