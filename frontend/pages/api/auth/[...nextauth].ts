import NextAuth from "next-auth"
import type { NextAuthOptions } from 'next-auth'
import GithubProvider from "next-auth/providers/github"

export const authOptions = {
  // Configure one or more authentication providers
  providers: [
    GithubProvider({
      clientId: process.env.GITHUB_ID as string,
      clientSecret: process.env.GITHUB_SECRET as string,
    }),
    // ...add more providers here
  ],
  jwt: {
    maxAge: 60
  },
  session: {
    maxAge: 60
  },
  callbacks: {
    async jwt({ token, account, profile }: any) {
      if (account) {
        console.log({ token, account, profile})
        token.accessToken = account.access_token
        token.id = profile.id
        token.providerAccountId = account.providerAccountId
        token.provider = account.provider
      }
      // console.log({ token })
      return token
    },
    async session({ session, token, user }: any) {
      // console.log({ session, token, user })
      return session;
    }
  }
}

export default NextAuth(authOptions as NextAuthOptions)