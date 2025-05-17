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

interface EmbeddingResponse {
	shape: number[];
	data: number[][];
}

interface PromptBody {
	data: string;
}

interface InsertBody {
	id: string;
	data: object;
	metadata: Record<string, VectorizeVectorMetadata>;
}

interface QueryBody {
	query: string;
}

export default {
	async fetch(request, env): Promise<Response> {


	async function readRequestBody<T>(request: Request) : Promise<T> {
		return await request.json();
	}


		let path = new URL(request.url).pathname;

		// Generate enhanced prompt from card data
		if (path === "/prompt") {
			const body = await readRequestBody<PromptBody>(request);

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

		if (path === "/insert") {
			const body = await readRequestBody<InsertBody>(request);

			// TODO: Consider parsing the payload into a different style string for better matching?
			const modelResp: EmbeddingResponse = await env.AI.run(
				"@cf/baai/bge-base-en-v1.5",
				{text: JSON.stringify(body.data)}
			);
	
			const vectors: VectorizeVector[] = modelResp.data.map(vector => ({
				id: body.id,
				values: vector,
				metadata: body.metadata
			}));
		
			let inserted = await env.VECTORIZE.upsert(vectors);
			return Response.json(inserted);
		}

		if (path === "/search") {
			const body = await readRequestBody<QueryBody>(request);

			var userQuery = body.query

			const queryVector: EmbeddingResponse = await env.AI.run(
				"@cf/baai/bge-base-en-v1.5",
				{
				text: [userQuery],
				},
			);
		
			const matches = await env.VECTORIZE.query(queryVector.data[0], {
				topK: 3,
				returnValues: true,
				returnMetadata: "all",
			});
			return Response.json({
				matches: matches,
			});
		}

		return Response.json({
			status: 404
		})
	},
} satisfies ExportedHandler<Env>;
