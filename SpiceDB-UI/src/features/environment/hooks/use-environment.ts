import { useQuery } from '@tanstack/react-query';
import { useNavigate, useParams } from '@tanstack/react-router'
import { getEnvNames } from '../api';

export interface SpiceDBEnvironment {
  name: string;
}

export function useEnvironment() {
  const navigate = useNavigate()
  const { envId } = useParams({ from: "/$envId" })

  const { data: envNames } = useQuery({
    queryKey: ['spicedb-env-names'],
    queryFn: getEnvNames,
  })
  const environments = envNames?.map(name => ({ name })) || [];

  const currentEnvironment = environments.find(e => e.name === envId) || null;

  const navigateToEnvironment = (envName: string) => {
    if (envName === envId) return;
    navigate({ to: "/$envId", params: { envId: envName } })
  };

  return {
    environments,
    currentEnvironment,
    navigateToEnvironment,
  };
}
