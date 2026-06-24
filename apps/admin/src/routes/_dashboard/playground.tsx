import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Send } from 'lucide-react'
import { useChat } from '@/hooks/use-chat'
import { ModelSelector } from '@/components/gateway/model-selector'
import { TokenCounter } from '@/components/gateway/token-counter'
import { useState } from 'react'

export const Route = createFileRoute('/_dashboard/playground')({
  component: PlaygroundPage,
})

function PlaygroundPage() {
  const [input, setInput] = useState('')
  const [model, setModel] = useState('mimo-auto')
  const { messages, sendMessage, isLoading } = useChat()

  const handleSubmit = () => {
    if (input.trim()) {
      sendMessage(input)
      setInput('')
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Playground</h1>
          <p className="text-muted-foreground">Test the MiMo gateway with interactive chat</p>
        </div>
        <ModelSelector value={model} onChange={setModel} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Input</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <Textarea
              placeholder="Type your message..."
              value={input}
              onChange={(e) => setInput(e.target.value)}
              rows={6}
            />
            <Button onClick={handleSubmit} disabled={isLoading || !input.trim()}>
              <Send className="h-4 w-4 mr-2" />
              {isLoading ? 'Sending...' : 'Send'}
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Response</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4 min-h-[200px]">
              {messages.length === 0 ? (
                <p className="text-muted-foreground text-center py-8">
                  No messages yet. Start a conversation!
                </p>
              ) : (
                messages.map((msg, i) => (
                  <div
                    key={i}
                    className={`p-3 rounded-lg ${
                      msg.role === 'user' ? 'bg-muted' : 'bg-accent/10'
                    }`}
                  >
                    <p className="text-sm font-medium mb-1">
                      {msg.role === 'user' ? 'You' : 'MiMo'}
                    </p>
                    <p className="text-sm">{msg.content}</p>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      <TokenCounter promptTokens={0} completionTokens={0} totalTokens={0} />
    </div>
  )
}
