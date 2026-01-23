export interface RouteStaticConfigI {
  header: HeaderConfigI;
}

export interface HeaderConfigI {
  isAuthenticated: boolean
}

export const RouteComponentTemplate: RouteStaticConfigI = {
  header: {
    isAuthenticated: false,
  }
}