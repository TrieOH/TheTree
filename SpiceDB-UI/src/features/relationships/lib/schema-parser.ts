export interface ParsedSchema {
  definitions: string[]
  relationsByDefinition: Record<string, string[]>
}

export function parseSpiceDBSchema(schemaText: string): ParsedSchema {
  const definitions: string[] = []
  const relationsByDefinition: Record<string, string[]> = {}

  // Basic regex to find definition blocks
  // Matches "definition <name> {" and "definition <name> {}"
  const definitionRegex = /definition\s+([a-zA-Z0-9_/]+)\s*{/g

  let match
  while ((match = definitionRegex.exec(schemaText)) !== null) {
    const defName = match[1]
    definitions.push(defName)

    // Find the end of this definition block to extract its relations
    const startIdx = match.index + match[0].length
    let endIdx = -1
    let braceCount = 1

    for (let i = startIdx; i < schemaText.length; i++) {
      if (schemaText[i] === '{') braceCount++
      if (schemaText[i] === '}') braceCount--
      if (braceCount === 0) {
        endIdx = i
        break
      }
    }

    if (endIdx !== -1) {
      const blockContent = schemaText.substring(startIdx, endIdx)

      // Match "relation <name>:"
      const relationRegex = /relation\s+([a-zA-Z0-9_]+)\s*:/g
      const relations: string[] = []
      let relMatch
      while ((relMatch = relationRegex.exec(blockContent)) !== null) {
        relations.push(relMatch[1])
      }

      // Match "permission <name> ="
      const permissionRegex = /permission\s+([a-zA-Z0-9_]+)\s*=/g
      while ((relMatch = permissionRegex.exec(blockContent)) !== null) {
        relations.push(relMatch[1])
      }

      relationsByDefinition[defName] = relations
    } else {
      relationsByDefinition[defName] = []
    }
  }

  return {
    definitions,
    relationsByDefinition,
  }
}
