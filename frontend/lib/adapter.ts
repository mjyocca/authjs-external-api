import { Adapter, AdapterUser } from 'next-auth/adapters';
import type { NextAuthOptions } from 'next-auth';
import { encode, JWT } from 'next-auth/jwt';
import fetch from './fetch';

const encodeJWT = async (token: JWT, maxAge: number) => {
  return await encode({
    token: {
      ...token,
      role: 'adapter',
    },
    secret: process.env.NEXTAUTH_SECRET as string,
    maxAge,
  });
};

export type AdapterConfig = Partial<NextAuthOptions> & {
  jwt: { maxAge: number };
};

export function MyAdapter(config: AdapterConfig): Adapter {
  let bearer: string = '';
  const getBearerToken = async () => {
    return await encodeJWT({}, config.jwt.maxAge);
  };

  const client = async (path: string, { headers, ...init }: RequestInit) => {
    if (!bearer) {
      bearer = await getBearerToken();
    }
    return await fetch(`${process.env.EXTERNAL_API_ENDPOINT}/${path}`, {
      headers: {
        'content-type': 'application/json',
        Authorization: `Bearer ${bearer}`,
        ...headers,
      },
      ...init,
    });
  };

  return {
    async createUser(user: Omit<AdapterUser, 'id'>) {
      const { name, email, image } = user;
      const result = await client(`api/adapter/user`, {
        method: 'POST',
        body: JSON.stringify({ name, email, image }),
      });
      if (!result.data) return user as AdapterUser;
      return { emailVerified: null, ...result.data } as AdapterUser;
    },
    async getUser(id) {
      const result = await client(`api/adapter/user?id=${encodeURIComponent(id)}`, {
        method: 'GET',
      });

      if (!result.data) return null;
      return result.data;
    },
    async getUserByEmail(email) {
      const result = await client(`api/adapter/user?email=${encodeURIComponent(email)}`, {
        method: 'GET',
      });

      if (!result.data) return null;
      return result.data;
    },
    async getUserByAccount({ providerAccountId, provider }) {
      const result = await client(
        `api/adapter/user?providerId=${encodeURIComponent(providerAccountId)}&providerType=${encodeURIComponent(
          provider
        )}`,
        {
          method: 'GET',
        }
      );
      if (!result.data) return null;
      return result.data;
    },
    async linkAccount(account) {
      const result = await client(`api/adapter/user`, {
        method: 'PATCH',
        body: JSON.stringify({ ...account }),
      });
    },
    // @ts-ignore
    createSession: () => {},
    // @ts-ignore
    getSessionAndUser: () => {},
    // @ts-ignore
    updateSession: () => {},
    // @ts-ignore
    deleteSession: () => {},
    // @ts-ignore
    updateUser: () => {},
    // @ts-ignore
    deleteUser: () => {},
  };
}
