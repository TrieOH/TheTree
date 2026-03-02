import { cn, maskStringMiddle } from '@/shared/lib/utils';
import { ShadowButton } from '@/shared/ui/buttons/ShadowButton';
import { Dialog, DialogClose, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/shared/ui/shadcn/dialog';
import { AlertCircle, Copy, Eye, EyeClosed } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { revokeApiKey, rotateApiKey } from '../api';

interface PropsI {
  publicKey: string;
}

export default function APIKeyManager ({ publicKey}: PropsI) {
  const [showPublic, setShowPublic] = useState(false);
  const [secretKey, setSecretKey] = useState<string | null>(null);
  const [showGenerateDialog, setShowGenerateDialog] = useState(false);
  const [showRevokeDialog, setShowRevokeDialog] = useState(false);

  const handlePublicKeyCopy = () => {
    navigator.clipboard.writeText(publicKey);
    toast.success("Project ID copied to clipboard");
  };

  const handleSecretKeyCopy = () => {
    if(secretKey) {
      navigator.clipboard.writeText(secretKey);
      toast.success("Secret Key copied to clipboard");
    } else toast.warning("Secret Key not generated");
  };

  const handleGenerateSecretKey = async () => {
    const response = await rotateApiKey(publicKey);
    if(response.success) {
      setShowGenerateDialog(false);
      setSecretKey(response.data.api_key);
      toast.success(response.message);
    }
  };

  const handleRevokeSecretKey = async () => {
    const response = await revokeApiKey(publicKey);
    if(response.success) {
      setShowRevokeDialog(false);
      setSecretKey(null);
      toast.success(response.message);
    }
  };

  return (
    <div className="max-w-2xl space-y-8">

      {/* Header */}
      <div className="space-y-1 border-b border-border pb-4 px-4">
        <h2 className="text-xl font-semibold text-primary">API Keys</h2>
        <p className="text-sm text-muted-foreground">Manage your API keys for accessing the platform.</p>
      </div>

      {/* Project ID Section */}
      <section className='bg-card py-4 rounded-sm border border-border'>
        <div className='border-b border-border pb-2 px-4'>
          <span>Project ID</span>
          <p className='text-muted-foreground text-sm'>
            Safe to embed in client-side code. This key never changes.
          </p>
        </div>
        <div
          className={cn(
            'mx-4 px-4 py-3 mt-4 rounded-sm bg-muted text-muted-foreground',
            'flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between'
          )}
        >
          <span className='break-all font-mono text-sm sm:text-base'>
            {showPublic ? publicKey : maskStringMiddle(publicKey, 4, 4)}
          </span>
          <div className='flex items-center justify-end gap-2 sm:justify-center'>
            <ShadowButton
              variant="ghost"
              onClick={() => setShowPublic(v => !v)}
              className='p-0'
              leftIcon={
                showPublic ? <EyeClosed className="h-4 w-4"/>
                : <Eye className="h-4.5 w-4.5"/>
              }
            />
            <ShadowButton
              variant="ghost"
              onClick={handlePublicKeyCopy}
              className='p-0'
              leftIcon={<Copy className="h-4 w-4"/>}
            />
          </div>
        </div>
      </section>

      {/* Secret Key Section */}
      <section className='bg-card py-4 rounded-sm border border-border'>
        <div className='border-b border-border pb-2 px-4'>
          <span>Secret Key</span>
          <p className='text-muted-foreground text-sm'>
            Only shown once upon generation. Store it securely — you cannot retrieve it again.
          </p>
        </div>
        {secretKey && (
          <div
            className={cn(
              'mx-4 px-4 py-3 mt-4 rounded-sm bg-muted text-muted-foreground',
              'flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between'
            )}
          >
            <span className='break-all font-mono text-sm sm:text-base'>
              {secretKey}
            </span>
            <div className='flex items-center justify-end gap-2 sm:justify-center'>
              <ShadowButton
                variant="ghost"
                onClick={handleSecretKeyCopy}
                className='p-0'
                leftIcon={<Copy className="h-4 w-4"/>}
              />
            </div>
          </div>
        )}
        <div className='mx-4 mt-4 flex flex-wrap justify-center gap-2'>
          <Dialog open={showGenerateDialog} onOpenChange={setShowGenerateDialog}>
            <DialogTrigger asChild>
              <ShadowButton className='h-8 px-3 text-xs' value='Generate Secret Key' variant="solid"/>
            </DialogTrigger>
            <DialogContent className="w-11/12 sm:max-w-lg">
              <DialogHeader>
                <DialogTitle className='flex sm:flex-row flex-col items-center gap-2'>
                  <AlertCircle className='h-5 w-5 text-destructive'/>
                  Are you absolutely sure?
                </DialogTitle>
                <DialogDescription>
                  This action cannot be undone. This will generate a new secret key.
                  Please ensure you save it in a secure location as it will only be shown once.
                </DialogDescription>
              </DialogHeader>
              <DialogFooter className="flex-col-reverse sm:flex-row sm:justify-end gap-2">
                <DialogClose asChild>
                  <ShadowButton variant="outline" className='w-full sm:w-auto' value='Cancel' />
                </DialogClose>
                <ShadowButton onClick={handleGenerateSecretKey} className='w-full sm:w-auto' value='Generate' />
              </DialogFooter>
            </DialogContent>
          </Dialog>

          <Dialog open={showRevokeDialog} onOpenChange={setShowRevokeDialog}>
            <DialogTrigger asChild>
              <ShadowButton variant="destructive" className='h-8 px-3 text-xs' value='Revoke Secret Key'/>
            </DialogTrigger>
            <DialogContent className="w-11/12 sm:max-w-lg">
              <DialogHeader>
                <DialogTitle className='flex sm:flex-row flex-col items-center gap-2'>
                  <AlertCircle className='h-5 w-5 text-destructive'/>
                  Are you absolutely sure?
                </DialogTitle>
                <DialogDescription>
                  This action cannot be undone. This will permanently revoke your current secret key,
                  and any applications using it will no longer be able to access the platform.
                </DialogDescription>
              </DialogHeader>
              <DialogFooter className="flex-col-reverse sm:flex-row sm:justify-end gap-2">
                <DialogClose asChild>
                  <ShadowButton variant="outline" className='w-full sm:w-auto' value='Cancel' />
                </DialogClose>
                <ShadowButton variant="destructive" onClick={handleRevokeSecretKey} className='w-full sm:w-auto' value='Revoke' />
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </div>
      </section>

    </div>
  );
};