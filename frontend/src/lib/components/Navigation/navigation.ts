

import { FileCode2, FolderClosed, Folders, House, Images, LayoutDashboard, LogOut, Settings, Users, type Icon as IconType } from "@lucide/svelte";
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
		href: `/admin/dashboard`,
		icon: LayoutDashboard
	},
	{
		label: 'Pages',
		href: `/admin/pages`,
		icon: FileCode2,
	},
	{
		label: 'Projects',
		href: `/admin/projects`,
		icon: FolderClosed,
	},
	{
		label: 'Media',
		href: `/admin/media`,
		icon: Images,
	},
	{
		label: 'Contacts',
		href: `/admin/contacts`,
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
