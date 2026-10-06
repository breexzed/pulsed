// Grounding plugin: forces source consultation before answering/editing.
const GROUNDED_TOOLS = new Set(["read", "grep", "glob", "webfetch", "websearch"])
const GATED_TOOLS = new Set(["edit", "write", "multiedit", "patch", "bash"])

const grounded: Record<string, boolean> = {}

const GROUNDING_RULE = [
  "GROUNDING RULE (mandatory):",
  "- Before answering any factual question about this codebase or the real world, consult a real source: read/grep the codebase, or websearch/webfetch the web.",
  "- Cite the source for every factual claim (file path, or URL).",
  "- If you cannot find a source, say so explicitly. Never answer from memory.",
  "- Prefer: codebase > project docs > official docs > general web.",
].join("\n")

const isQuestion = (text: string) =>
  /\?/.test(text) || /^(why|how|what|when|where|which|who|explain|does|is|can|should|would|compare|difference)\b/i.test(text.trim())

export const GroundingPlugin = async () => {
  return {
    "experimental.chat.system.transform": async (_input: any, output: any) => {
      output.system.push(GROUNDING_RULE)
    },

    "chat.message": async (input: any, output: any) => {
      grounded[input.sessionID] = false
      const text = output.parts
        .filter((p: any) => p.type === "text")
        .map((p: any) => p.text)
        .join(" ")
      const part = output.parts.find((p: any) => p.type === "text")
      if (isQuestion(text) && part) {
        part.text += "\n\n[GROUNDING REQUIRED: consult the codebase (read/grep) or the web (websearch/webfetch) before answering. Cite sources. If unavailable, say so.]"
      }
    },

    "tool.execute.before": async (input: any, _output: any) => {
      if (GATED_TOOLS.has(input.tool) && !grounded[input.sessionID]) {
        throw new Error(
          `Blocked: '${input.tool}' requires grounding first this turn. Call read/grep/websearch/webfetch on the relevant source, then retry.`,
        )
      }
    },

    "tool.execute.after": async (input: any, _output: any) => {
      if (GROUNDED_TOOLS.has(input.tool)) {
        grounded[input.sessionID] = true
      }
    },
  }
}

export default GroundingPlugin
