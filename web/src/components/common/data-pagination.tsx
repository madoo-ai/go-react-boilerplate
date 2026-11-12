import { useTranslations } from "@/lib/i18n"
import { Button } from "@/components/ui/button"
import { Pagination, PaginationContent, PaginationEllipsis, PaginationItem } from "@/components/ui/pagination"

function pageWindow(page: number, pageCount: number) {
  const pages = new Set<number>([1, pageCount, page])
  if (page > 1) pages.add(page - 1)
  if (page < pageCount) pages.add(page + 1)
  return [...pages].sort((a, b) => a - b)
}

export function DataPagination({
  page,
  pageCount,
  total,
  pageSize,
  onPageChange,
}: {
  page: number
  pageCount: number
  total: number
  pageSize: number
  onPageChange: (page: number) => void
}) {
  const t = useTranslations("common")

  if (pageCount <= 1) {
    return (
      <p className="text-sm text-muted-foreground">
        {t("showing", { from: total === 0 ? 0 : 1, to: Math.min(pageSize, total), total })}
      </p>
    )
  }

  const pages = pageWindow(page, pageCount)
  const from = (page - 1) * pageSize + 1
  const to = Math.min(page * pageSize, total)

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <p className="text-sm text-muted-foreground">
        {t("showing", { from, to, total })}
      </p>
      <Pagination className="mx-0 w-auto justify-start sm:justify-end">
        <PaginationContent>
          <PaginationItem>
            <Button variant="ghost" size="sm" disabled={page <= 1} onClick={() => onPageChange(page - 1)}>
              {t("previous")}
            </Button>
          </PaginationItem>
          {pages.map((pageNumber, index) => {
            const previous = pages[index - 1]
            const hasGap = previous !== undefined && pageNumber - previous > 1
            return (
              <PaginationItem key={pageNumber}>
                {hasGap ? <PaginationEllipsis /> : null}
                <Button
                  variant={pageNumber === page ? "outline" : "ghost"}
                  size="icon"
                  onClick={() => onPageChange(pageNumber)}
                >
                  {pageNumber}
                </Button>
              </PaginationItem>
            )
          })}
          <PaginationItem>
            <Button variant="ghost" size="sm" disabled={page >= pageCount} onClick={() => onPageChange(page + 1)}>
              {t("next")}
            </Button>
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  )
}
