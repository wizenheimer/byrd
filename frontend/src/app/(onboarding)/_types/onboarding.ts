export const DEFAULT_PROFILE = [
	"branding",
	"customers",
	"integration",
	"product",
	"pricing",
	"partnerships",
	"messaging",
] as const;

export type ProfileType = (typeof DEFAULT_PROFILE)[number];

// API Request/Response Types
export interface WorkspaceCreationRequest {
	competitors: string[];
	profiles: ProfileType[];
	features: string[];
}
