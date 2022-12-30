import NextAuth from 'next-auth';
import type { NextAuthOptions } from 'next-auth';
import GithubProvider from 'next-auth/providers/github';
import { MyAdapter } from '../../../lib/adapter';
import type { AdapterConfig } from '../../../lib/adapter';

const MAX_AGE = 60 * 60;

const jwtOptions: Partial<NextAuthOptions['jwt']> = {
  maxAge: MAX_AGE,
};

const { log } = console;

export const authOptions: NextAuthOptions = {
  providers: [
    GithubProvider({
      clientId: process.env.GITHUB_ID as string,
      clientSecret: process.env.GITHUB_SECRET as string,
    }),
  ],
  session: {
    strategy: 'jwt',
    maxAge: jwtOptions.maxAge,
  },
  jwt: jwtOptions,
  callbacks: {
    jwt: async ({ token, user, account, profile, isNewUser }) => {
      log('callback.jwt', { token, user, account, profile, isNewUser });
      if (account) {
        log('callback.jwt.account');
        token.accessToken = account.access_token;
        token.providerAccountId = account.providerAccountId;
        token.provider = account.provider;
      }
      if (user) {
        token.id = user.id;
      }
      return token;
    },
    async session({ session, token, user }: any) {
      log('callback.session', { session, token, user });
      return session;
    },
  },
};

/* initialize Custom Adapter */
authOptions.adapter = MyAdapter(authOptions as AdapterConfig);

export default NextAuth(authOptions as NextAuthOptions);
