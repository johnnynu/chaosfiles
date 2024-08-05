import { atom } from "jotai";
import { AuthUser } from "aws-amplify/auth";

export const userAtom = atom<AuthUser | null>(null);
export const errorAtom = atom<unknown>(null);
