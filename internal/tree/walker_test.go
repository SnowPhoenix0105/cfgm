package tree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalker_ImplementInterface(t *testing.T) {
	var _ NodeReadWriter = &walker{}
}

func TestWalker_EnterExit(t *testing.T) {
	root := NewNode()
	buildTime := ModifyTime(1)
	modifyTime := ModifyTime(2)
	build := func() {
		handler := WriteFrom(root, buildTime)
		assert.Equal(t, buildTime, handler.(*walker).time)

		handler.EnterObj("Server")
		handler.SetBool(true)
		func() {
			handler.EnterObj("Port")
			handler.SetInt(8080)
			handler.Exit()

			handler.EnterObj("Host")
			handler.SetString("localhost")
			handler.Exit()
		}()
		handler.Exit()
		assert.Empty(t, handler.(*walker).stack)
		assert.Equal(t, root, handler.(*walker).currentNode)
		assert.Equal(t, buildTime, root.ModifyTime())
		assert.True(t, root.Has(NodeKeyObj))
	}
	modify := func() {
		handler := WriteFrom(root, modifyTime)
		assert.Equal(t, modifyTime, handler.(*walker).time)

		handler.EnterObj("Server")
		func() {
			handler.EnterObj("CORS")
			func() {
				handler.EnterList(1)
				handler.SetString("localhost:8001")
				handler.Exit()
				handler.EnterList(0)
				handler.SetString("localhost:8000")
				handler.Exit()
				handler.EnterList(4)
				handler.SetString("localhost:8004")
				handler.Exit()
				handler.EnterList(3)
				handler.SetString("localhost:8003")
				handler.Exit()
				handler.EnterList(2)
				handler.SetString("localhost:8002")
				handler.Exit()
			}()
			handler.Exit()
		}()
		handler.Exit()
		assert.Equal(t, modifyTime, root.ModifyTime())
	}
	read := func() {
		handler := ReadFrom(root)

		assert.True(t, handler.TryEnterObj("Server"))
		assert.True(t, handler.Bool())
		assert.Equal(t, modifyTime, handler.ModifyTime())
		func() {

			assert.True(t, handler.TryEnterObj("Port"))
			assert.Equal(t, int64(8080), handler.Int())
			assert.Equal(t, buildTime, handler.ModifyTime())
			handler.Exit()

			assert.True(t, handler.TryEnterObj("Host"))
			assert.Equal(t, "localhost", handler.String())
			assert.Equal(t, buildTime, handler.ModifyTime())
			handler.Exit()

			assert.True(t, handler.TryEnterObj("CORS"))
			assert.Equal(t, modifyTime, handler.ModifyTime())
			func() {
				for i := 0; i <= 4; i++ {
					assert.True(t, handler.TryEnterList(i))
					assert.Equal(t, fmt.Sprintf("localhost:800%d", i), handler.String())
					assert.Equal(t, modifyTime, handler.ModifyTime())
					handler.Exit()
				}
			}()
			handler.Exit()
		}()
		handler.Exit()
	}
	build()
	modify()
	read()
}
