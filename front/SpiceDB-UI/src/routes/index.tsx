import { redirect, createFileRoute } from '@tanstack/react-router'
import { getEnvNames } from '#/features/environment/api';

export const Route = createFileRoute('/')({
  loader: async () => {
    const environments = await getEnvNames();
    if (environments.length === 0) throw new Error('No SpiceDB environments defined.')
    return redirect({ to: "/$envId", params: { envId: environments[0] } })
  },
})
