import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({ 
  staticData: {
    components: {
      header: {
        isAuthenticated: true,
      }
    }
  },
  component: App
})

function App() {

  return (
    <div className="bg-background">

    </div>
  )
}
