import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import type { StoreApi, UseBoundStore } from "zustand";

type WithSelectors<S> = S extends { getState: () => infer T }
	? S & { use: { [K in keyof T]: () => T[K] } }
	: never;

export const createSelectors = <S extends UseBoundStore<StoreApi<object>>>(
	_store: S,
) => {
	const store = _store as WithSelectors<typeof _store>;
	store.use = {};
	for (const k of Object.keys(store.getState())) {
		// biome-ignore lint/suspicious/noExplicitAny: ignored because it's a type assertion that's safe
		(store.use as any)[k] = () => store((s) => s[k as keyof typeof s]);
	}
	return store;
};

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}
