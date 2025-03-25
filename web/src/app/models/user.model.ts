export interface User {
  id?: number;
  name: string;
  email: string;
  isSuper: boolean;
  accessToken?: string;
  refreshToken?: string;
  isAuthenticated: boolean;
}
