export const ssr = true
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
  const res = await fetch('http://localhost:8080/api/v1/home');
  const json = await res.json();



  const homeGallery = json.data.galleries.find((item: any) => item.gallery.slug === 'home-hero-gallery');
  console.log('Found homeGallery:', homeGallery);

  return {
    Page: json.data.page,
    Content: json.data.page.content,
    HomeGallery: homeGallery,
    Galleries: json.data.galleries
  };
};
