/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run `npm run dev` in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run `npm run deploy` to publish your worker
 *
 * Bind resources to your worker in `wrangler.jsonc`. After adding bindings, a type definition for the
 * `Env` object can be regenerated with `npm run cf-typegen`.
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

const SEMANTIC_METADATA_PROMPT = "In a phrase, semantically and accurately describe the effects and behaviours of this Magic: The Gathering card. Omit proper grammar in favour of a shorter response with clear meaning for a text-embedded vector db. Do not wrap it in speechmarks."

export interface Env {
	VECTORIZE: Vectorize;
	AI: Ai;
}

export interface Payload {
	data: string;
}

export default {
	async fetch(request, env): Promise<Response> {


	async function readRequestBody(request: Request) : Promise<Payload> {
		return await request.json();
	}

		let path = new URL(request.url).pathname;

		// Generate enhanced prompt from card data
		if (path === "/prompt") {
			const body = await readRequestBody(request);

			const messages = [
				{ role: "system", content: "You are a Magic: The Gathering player who is prone to using both generic and magic-specific terminology" },
				{
					role: "user",
					content: SEMANTIC_METADATA_PROMPT + JSON.stringify(body).replace(/\"/g, "'"),
				},
			];

			const response = await env.AI.run("@cf/meta/llama-4-scout-17b-16e-instruct", { messages });

			return Response.json(response);
		}
	},
} satisfies ExportedHandler<Env>;
