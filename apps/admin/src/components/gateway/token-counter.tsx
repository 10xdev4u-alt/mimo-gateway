interface TokenCounterProps {
  promptTokens: number
  completionTokens: number
  totalTokens: number
}

export function TokenCounter({ promptTokens, completionTokens, totalTokens }: TokenCounterProps) {
  return (
    <div className="flex items-center gap-4 text-sm text-muted-foreground">
      <span>Prompt: {promptTokens}</span>
      <span>Completion: {completionTokens}</span>
      <span className="font-medium text-foreground">Total: {totalTokens}</span>
    </div>
  )
}
