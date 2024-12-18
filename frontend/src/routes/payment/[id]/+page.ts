import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';
import type { Payment } from '$lib/types';

export const load: PageLoad = async({ params, fetch }) => {
  const { id } = params;

    const response = await fetch(`/api/payment/${id}`);
    if (response.ok) {
      const payment = await response.json() as Payment;
      return { payment };
    } else {
      return error(404, "Payment not found");
    }
};