
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
  const res = await fetch('http://localhost:8080/api/v1/home'); // <-- use actual backend API address
  const json = await res.json();


  console.log('pageserverload data:', json.data)

  return {
    Hero: `This is where the Hero Element with data ${json.data} belongs`
  };
};

