import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { requireAuth } from '@/hooks/auth-guard'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/projects/')({
  beforeLoad: requireAuth,
  staticData: {
    components: {
      header: {
        test: false,
      }
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  // console.log(env.VITE_API_URL)
  return (
    <main className="w-full bg-background">
      <div>
        <Tabs defaultValue="projects" className='mt-6'>
          <TabsList className='w-full justify-start gap-4 bg-transparent border-b rounded-none py-0 px-4'>
            <TabsTrigger
              value="projects"
              className='max-w-60 h-full shadow-none! border-b data-[state=active]:border-b-black rounded-none'
            >
              Projetos
            </TabsTrigger>
            <TabsTrigger
              value="settings"
              className='max-w-60 h-full shadow-none! border-b data-[state=active]:border-b-black rounded-none'
            >
              Configurações
            </TabsTrigger>
          </TabsList>
          <TabsContent className='p-12' value="projects">What</TabsContent>
          <TabsContent 
            className='p-12' 
            value="settings"
          >
            Make changes to your account here.
          </TabsContent>
        </Tabs>
      </div>
    </main>
  )
}
