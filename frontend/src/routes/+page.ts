import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = () => {
  // This works on both server and client side
  // Status 307 is a temporary redirect; 308 is permanent
  throw redirect(307, '/admin/dashboard');
};
