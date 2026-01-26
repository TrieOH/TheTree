export interface RouteStaticConfigI {
  header: HeaderConfigI;
}

export interface HeaderConfigI {
  test: boolean
}

export const RouteComponentTemplate: RouteStaticConfigI = {
  header: {
    test: false,
  }
}