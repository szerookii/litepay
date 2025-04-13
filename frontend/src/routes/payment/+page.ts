import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import type { Payment } from "$lib/types";
import { z } from "zod";

const schema = z.object({
	id: z.string().uuid().nonempty()
});

export const load: PageLoad = async ({ url, fetch }) => {
	const data = schema.safeParse(Object.fromEntries(new URLSearchParams(url.search)));

	if (!data.success) {
		return error(400, "Invalid payment ID");
	}

	const { id } = data.data;

	const response = await fetch(`/api/payment/${id}`);
	if (response.ok) {
		const payment = (await response.json()) as Payment;
		return { payment };
	} else {
		return error(404, "Payment not found");
	}
};
