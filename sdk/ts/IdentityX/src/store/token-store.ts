let _accessToken: string | null = null;

export const tokenStore = {
  getAccessToken: () => _accessToken,
  setAccessToken: (token: string | null) => {
    _accessToken = token;
  },
  clear: () => {
    _accessToken = null;
  }
};
