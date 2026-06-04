import { useMutation } from '@tanstack/react-query'
import {
  openFormOnNamespaceFn,
  closeFormOnNamespaceFn,
  archiveFormOnNamespaceFn,
  redraftFormOnNamespaceFn,
} from '#/features/namespaces/api'
import {
  openFormFn,
  closeFormFn,
  archiveFormFn,
  redraftFormFn,
} from '#/features/forms/api'
import { Button } from '#/shared/ui/shadcn/button'
import { toast } from 'sonner'
import {
  FormStatusArchived,
  FormStatusClosed,
  FormStatusDraft,
  FormStatusOpen
} from "@trieoh/informd-models"
import { cn } from '#/shared/lib/utils'
import {
  Archive,
  Play,
  RotateCcw,
  StopCircle,
} from 'lucide-react'
import type { FormI } from '#/features/forms/model'
import { useState } from 'react'
import { ConfirmModal } from '#/widgets/modal/modal'

interface FormAdminHeaderProps {
  title: string
  description: string
  form: FormI
  namespaceID?: string
  responseCount: number
  onUpdate: (updatedForm: FormI) => void
}

export default function FormAdminHeader({
  title,
  description,
  form,
  namespaceID,
  responseCount,
  onUpdate
}: FormAdminHeaderProps) {
  const [isArchiveModalOpen, setIsArchiveModalOpen] = useState(false)

  const statusConfig = {
    [FormStatusOpen]: { label: 'Open', color: 'bg-green-500' },
    [FormStatusDraft]: { label: 'Draft', color: 'bg-slate-400' },
    [FormStatusClosed]: { label: 'Closed', color: 'bg-red-500' },
    [FormStatusArchived]: { label: 'Archived', color: 'bg-amber-500' },
  }
  const currentStatus = statusConfig[form.status]

  const { mutate: openForm, isPending: isOpenPending } = useMutation({
    mutationFn: () => namespaceID ? openFormOnNamespaceFn(namespaceID, form.id) : openFormFn(form.id),
    onSuccess: (response) => {
      if (response.success) {
        onUpdate(response.data)
        toast.success('Form opened successfully')
      } else toast.error(response.message || 'Failed to open form')
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const { mutate: closeForm, isPending: isClosePending } = useMutation({
    mutationFn: () => namespaceID ? closeFormOnNamespaceFn(namespaceID, form.id) : closeFormFn(form.id),
    onSuccess: (response) => {
      if (response.success) {
        onUpdate(response.data)
        toast.success('Form closed successfully')
      } else toast.error(response.message || 'Failed to close form')
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const { mutate: archiveForm, isPending: isArchivePending } = useMutation({
    mutationFn: () => namespaceID ? archiveFormOnNamespaceFn(namespaceID, form.id) : archiveFormFn(form.id),
    onSuccess: (response) => {
      if (response.success) {
        onUpdate(response.data)
        toast.success('Form archived successfully')
        setIsArchiveModalOpen(false)
      } else toast.error(response.message || 'Failed to archive form')
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const { mutate: redraftForm, isPending: isRedraftPending } = useMutation({
    mutationFn: () => namespaceID ? redraftFormOnNamespaceFn(namespaceID, form.id) : redraftFormFn(form.id),
    onSuccess: (response) => {
      if (response.success) {
        onUpdate(response.data)
        toast.success('Form redrafted successfully')
      } else toast.error(response.message || 'Failed to redraft form')
    },
    onError: (error: Error) => toast.error(error.message),
  })

  const isPending = isOpenPending || isClosePending || isArchivePending || isRedraftPending

  return (
    <div className="flex flex-col gap-4 w-full">
      <div className="flex flex-col gap-0.5">
        <h1 className="text-lg font-semibold tracking-tight wrap-break-word">{title}</h1>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>

      <div className="flex flex-wrap items-center gap-x-6 gap-y-4">
        {/* Status Badge Group */}
        <div className="flex items-center gap-2 pr-6 border-r border-border/60 last:border-r-0">
          <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Status</span>
          <div className="flex items-center gap-1.5 px-2 py-1 bg-secondary/40 rounded-sm border border-border/40">
            <div className={cn("size-1.5 rounded-full", currentStatus.color)} />
            <span className="text-[10px] font-bold uppercase tracking-widest">{currentStatus.label}</span>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-wrap items-center gap-2">
          {form.status === FormStatusDraft && (
            <Button
              size="sm" onClick={() => openForm()}
              disabled={isPending}
              className="h-8 rounded-sm text-[10px] font-bold uppercase tracking-wider"
            >
              <Play className="mr-1.5 size-3" />
              Open
            </Button>
          )}

          {form.status === FormStatusOpen && (
            <>
              <Button
                size="sm"
                variant="outline"
                onClick={() => redraftForm()}
                disabled={isPending || responseCount > 0}
                className="h-8 rounded-sm text-[10px] font-bold uppercase tracking-wider"
              >
                <RotateCcw className="mr-1.5 size-3" />
                Draft
              </Button>
              <Button
                size="sm" variant="destructive"
                onClick={() => closeForm()}
                disabled={isPending}
                className="h-8 rounded-sm text-[10px] font-bold uppercase tracking-wider"
              >
                <StopCircle className="mr-1.5 size-3" />
                Close
              </Button>
            </>
          )}

          {form.status === FormStatusClosed && (
            <>
              <Button
                size="sm"
                variant="outline"
                onClick={() => setIsArchiveModalOpen(true)}
                disabled={isPending}
                className="h-8 rounded-sm text-[10px] font-bold uppercase tracking-wider"
              >
                <Archive className="mr-1.5 size-3" />
                Archive
              </Button>
              <ConfirmModal
                isOpen={isArchiveModalOpen}
                onClose={() => setIsArchiveModalOpen(false)}
                onConfirm={() => archiveForm()}
                title="Archive Form"
                description="Are you sure you want to archive this form? This action cannot be undone and the form will no longer be accessible for responses."
                confirmText="Archive"
                variant="destructive"
                isLoading={isArchivePending}
              />
            </>
          )}
        </div>
      </div>
    </div>
  )
}

