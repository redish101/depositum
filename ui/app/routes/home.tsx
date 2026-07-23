import { useGetObjects } from "~/lib/api/generated/object/object"

import { Button } from "~/components/ui/button"
import { Badge } from "~/components/ui/badge"
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table"

export default function Home() {
  const { data, error, isLoading, mutate } = useGetObjects()

  const response = data?.data

  const pagination = response && "items" in response ? response : undefined

  const objects = pagination?.items ?? []

  return (
    <div className="flex flex-1 flex-col">
      <div className="@container/main flex flex-1 flex-col gap-6 p-4 md:p-6">
        {/* Header */}

        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">Objects</h1>

            <p className="text-sm text-muted-foreground">
              管理档案对象和同步状态
            </p>
          </div>

          <Button variant="outline" onClick={() => mutate()}>
            Refresh
          </Button>
        </div>

        {/* Statistics */}

        <div className="grid gap-4 md:grid-cols-3">
          <StatCard title="Total Objects" value={pagination?.total} />

          <StatCard
            title="API Status"
            value={isLoading ? "Loading" : error ? "Error" : "Connected"}
          />

          <StatCard title="Current Page" value={pagination?.page} />
        </div>

        {/* Table */}

        <Card>
          <CardHeader>
            <CardTitle>Object List</CardTitle>
          </CardHeader>

          <CardContent>
            {isLoading && (
              <div className="py-8 text-center text-sm text-muted-foreground">
                Loading objects...
              </div>
            )}

            {error && (
              <div className="rounded-lg bg-destructive/10 p-4 text-sm text-destructive">
                Failed to load objects
              </div>
            )}

            {!isLoading && !error && (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Name</TableHead>

                    <TableHead>Library</TableHead>

                    <TableHead>Shelf</TableHead>

                    <TableHead>Status</TableHead>
                  </TableRow>
                </TableHeader>

                <TableBody>
                  {objects.length === 0 && (
                    <TableRow>
                      <TableCell
                        colSpan={4}
                        className="h-24 text-center text-muted-foreground"
                      >
                        No objects found
                      </TableCell>
                    </TableRow>
                  )}

                  {objects.map((object) => (
                    <TableRow key={object.id}>
                      <TableCell className="font-medium">
                        {object.name}
                      </TableCell>

                      <TableCell>
                        {object.currentStatus?.library?.name ?? "-"}
                      </TableCell>

                      <TableCell>
                        {object.currentStatus?.shelf?.name ?? "-"}
                      </TableCell>

                      <TableCell>
                        <Badge variant="outline">
                          {object.currentStatus?.phase ?? "UNKNOWN"}
                        </Badge>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

function StatCard({
  title,
  value,
}: {
  title: string
  value?: string | number
}) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
      </CardHeader>

      <CardContent>
        <div className="text-2xl font-bold">{value ?? "-"}</div>
      </CardContent>
    </Card>
  )
}
