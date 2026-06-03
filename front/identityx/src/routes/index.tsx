import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  staticData: {
    components: {
      header: "landing"
    }
  },
  component: App
})

function App() {

  return (
    <div className="bg-background min-h-screen flex items-center justify-center">
      <div className="bg-white max-w-lg p-8 rounded-2xl shadow-lg">

        <h1 className="text-3xl font-bold text-center mb-6">
          🍰 Bolo
        </h1>

        <h2 className="text-xl font-semibold mb-2">
          Ingredientes
        </h2>

        <ul className="list-disc pl-5 space-y-1 text-gray-700">
          <li>2 xícaras de farinha de trigo</li>
          <li>1 xícara de açúcar</li>
          <li>1 xícara de leite</li>
          <li>3 ovos</li>
          <li>3 colheres de manteiga</li>
          <li>1 colher de sopa de fermento</li>
        </ul>

        <h2 className="text-xl font-semibold mt-6 mb-2">
          Modo de preparo
        </h2>

        <ol className="list-decimal pl-5 space-y-1 text-gray-700">
          <li>Preaqueça o forno a 180°C</li>
          <li>Misture ovos, açúcar e manteiga</li>
          <li>Adicione o leite e a farinha aos poucos</li>
          <li>Misture até ficar homogêneo</li>
          <li>Adicione o fermento e mexa levemente</li>
          <li>Coloque em forma untada</li>
          <li>Asse por cerca de 35 minutos</li>
        </ol>

      </div>
    </div>
  )
}
