package runwrappers

import (
	"context"
	"testing"

	"github.com/spf13/cobra"
)

func TestChain(t *testing.T) {
	t.Parallel()

	newRunWithCtx := func(ch chan int, val int) RunEWithContext {
		return func(ctx context.Context, cmd *cobra.Command, args []string) error {
			ch <- val
			return nil
		}
	}

	mockWrap := func(wrap RunEWithContext) WrapRunEWithContext {
		return func(next RunEWithContext) RunEWithContext {
			return func(ctx context.Context, cmd *cobra.Command, args []string) error {
				wrap(ctx, cmd, args)
				return next(ctx, cmd, args)
			}
		}
	}

	t.Run("ok with multiple wrappers", func(t *testing.T) {
		t.Parallel()

		ch := make(chan int, 5)

		chain := Chain(
			mockWrap(newRunWithCtx(ch, 1)),
			mockWrap(newRunWithCtx(ch, 2)),
			mockWrap(newRunWithCtx(ch, 3)),
			mockWrap(newRunWithCtx(ch, 4)),
		)

		chain(newRunWithCtx(ch, 5))(context.Background(), nil, nil)

		for i := 1; i <= 5; i++ {
			got := <-ch
			if got != i {
				t.Errorf("for chained wrapper call ordering got call number %v in call number %v", got, i)
			}
		}
	})

	t.Run("ok with one wrapper", func(t *testing.T) {
		t.Parallel()

		ch := make(chan int, 2)

		chain := Chain(mockWrap(newRunWithCtx(ch, 1)))

		chain(newRunWithCtx(ch, 2))(context.Background(), nil, nil)

		for i := 1; i <= 2; i++ {
			got := <-ch
			if got != i {
				t.Errorf("for chained wrapper call ordering got call number %v in call number %v", got, i)
			}
		}
	})
}
