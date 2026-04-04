export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  plan: string;
  monthly_limit: number;
}

export interface OrgListItem {
  id: string;
  name: string;
  slug: string;
  role: string;
}

export interface MeResponse {
  user: User;
  organization: Organization;
  role: string;
  organizations: OrgListItem[];
}
