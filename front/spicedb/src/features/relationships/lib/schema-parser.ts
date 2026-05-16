export interface RelationInfo {
  name: string
  allowedSubjectTypes: string[]
  type: 'relation' | 'permission'
}

export interface ParsedSchema {
  definitions: string[]
  relationsByDefinition: Record<string, RelationInfo[] | undefined>
}

export function parseSpiceDBSchema(schemaText: string): ParsedSchema {
  const definitions: string[] = []
  const relationsByDefinition: Record<string, RelationInfo[] | undefined> = {}

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

    const relations: RelationInfo[] = []
    if (endIdx !== -1) {
      const blockContent = schemaText.substring(startIdx, endIdx)

      const lines = blockContent.split('\n')
      for (let line of lines) {
        line = line.trim()
        if (line.startsWith('relation ')) {
          // Match "relation <name>: <types>"
          const relMatch = line.match(/relation\s+([a-zA-Z0-9_]+)\s*:\s*([^;]+)/)
          if (relMatch) {
            const name = relMatch[1]
            const typesStr = relMatch[2].split('//')[0] // remove comments
            const allowedSubjectTypes = Array.from(
              new Set(
                typesStr
                  .split('|')
                  .map((t) => t.trim())
                  .map((t) => {
                    // Handle type#rel or type:* or type
                    return t.split('#')[0].split(':')[0].trim()
                  })
                  .filter((t) => t !== ''),
              ),
            )

            relations.push({
              name,
              allowedSubjectTypes,
              type: 'relation',
            })
          }
        } else if (line.startsWith('permission ')) {
          // Match "permission <name> ="
          const permMatch = line.match(/permission\s+([a-zA-Z0-9_]+)\s*=\s*(.+)/)
          if (permMatch) {
            relations.push({
              name: permMatch[1],
              allowedSubjectTypes: [], // Permissions don't explicitly list allowed types
              type: 'permission',
            })
          }
        }
      }
    }
    relationsByDefinition[defName] = relations
  }

  return {
    definitions,
    relationsByDefinition,
  }
}
