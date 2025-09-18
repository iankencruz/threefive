/* eslint-disable @typescript-eslint/no-explicit-any */
export const ssr = true;
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
  const res = await fetch('http://localhost:8080/api/v1/home');
  const json = await res.json();

  console.log('data', json);
  // console.log('hero_gallery: ', json.data.hero_gallery);



  return {
    Page: json.data.page,
    Content: json.data.content,
    HeroGallery: json.data.hero_gallery.media
  };
};
