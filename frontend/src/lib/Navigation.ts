
import { FileCode2, FolderClosed, Folders, House, LogOut, Settings, Users, type Icon as IconType } from "@lucide/svelte";
import { base } from '$app/paths';
export interface NavigationItem {
	label: string;
	href?: string;
	icon?: typeof IconType, // optional icon component
	children?: NavigationItem[];
	permissions?: string[]; // optional role/permission gate
	activeMatch?: string;   // optional override for active route matching
	action?: string; // 'logout' or other custom action
}


export const userMenuItems: NavigationItem[] = [
	{ label: 'Settings', href: '/settings', icon: Settings },
	{ label: 'Logout', action: 'logout', icon: LogOut }
];


export const sidebarNavigation: NavigationItem[] = [
	{
		label: 'Dashboard',
		href: `${base}/dashboard`,
		icon: House
	},
	{
		label: 'Pages',
		href: `${base}/pages`,
		icon: FileCode2,
	},
	{
		label: 'Projects',
		href: `${base}/projects`,
		icon: FolderClosed,
	},
	{
		label: 'Contacts',
		href: `${base}/contacts`,
		icon: Users
	},
	// {
	// 	label: 'Invoices',
	// 	href: '/invoices',
	// 	icon: 'receipt'
	// },
	// {
	// 	label: 'Settings',
	// 	href: '/settings',
	// 	icon: 'cog'
	// }
];

