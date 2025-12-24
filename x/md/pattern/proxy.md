结构型模式

代理模式（Proxy Pattern）
作用：为其他对象提供代理以控制对其的访问（如延迟加载、权限控制）。

Go 实现：通过包装对象并实现相同接口。

```go
type Image interface {
    Display()
}

type RealImage struct{ filename string }
func (r *RealImage) Display() { fmt.Println("Displaying", r.filename) }

type ProxyImage struct {
    realImage *RealImage
    filename  string
}
func (p *ProxyImage) Display() {
    if p.realImage == nil {
        p.realImage = &RealImage{p.filename} // 延迟加载
    }
    p.realImage.Display()
}

// 使用
image := &ProxyImage{filename: "test.jpg"}
image.Display() // 首次调用时加载真实对象
```

代理模式是一种结构型设计模式，它允许通过代理对象来控制对另一个对象（目标对象）的访问。代理对象充当了目标对象的接口，客户端通过代理对象来间接访问目标对象。
在 Go 语言中，代理模式可以通过接口和结构体来实现。下面为你提供几种不同场景下的代理模式示例。

远程代理
远程代理用于为一个位于不同地址空间的对象提供一个本地的代理对象，客户端通过这个代理对象来访问远程对象。这里模拟一个简单的远程服务调用：
```go
package main

import (
    "fmt"
)

// Service 定义服务接口
type Service interface {
    Request() string
}

// RealService 实现服务接口，代表远程服务
type RealService struct{}

func (rs *RealService) Request() string {
    return "RealService: Handling request"
}

// RemoteProxy 远程代理
type RemoteProxy struct {
    realService *RealService
}

func NewRemoteProxy() *RemoteProxy {
    return &RemoteProxy{
       realService: &RealService{},
    }
}

func (rp *RemoteProxy) Request() string {
    // 模拟远程调用的额外操作，如网络连接等
    fmt.Println("RemoteProxy: Connecting to remote service...")
    result := rp.realService.Request()
    fmt.Println("RemoteProxy: Disconnecting from remote service...")
    return result
}

func main() {
    proxy := NewRemoteProxy()
    response := proxy.Request()
    fmt.Println(response)
}

```
代码解释：

Service 是服务接口，定义了 Request 方法。
RealService 是实际的服务实现，代表远程服务。
RemoteProxy 是远程代理，它持有 RealService 的实例。在 Request 方法中，模拟了远程调用的额外操作，如连接和断开远程服务。
在 main 函数中，客户端通过代理对象调用服务方法。

虚拟代理
虚拟代理用于在需要创建开销很大的对象时，延迟对象的创建，直到真正需要使用该对象时才进行创建。以下是一个简单的图片加载的虚拟代理示例：
```go
package main

import (
    "fmt"
)

// Image 定义图片接口
type Image interface {
    Display()
}

// RealImage 实现图片接口，代表真实的图片对象
type RealImage struct {
    filename string
}

func NewRealImage(filename string) *RealImage {
    img := &RealImage{filename: filename}
    img.loadFromDisk()
    return img
}

func (ri *RealImage) loadFromDisk() {
    fmt.Printf("Loading image: %s\n", ri.filename)
}

func (ri *RealImage) Display() {
    fmt.Printf("Displaying image: %s\n", ri.filename)
}

// ProxyImage 虚拟代理
type ProxyImage struct {
    filename  string
    realImage *RealImage
}

func NewProxyImage(filename string) *ProxyImage {
    return &ProxyImage{
       filename: filename,
    }
}

func (pi *ProxyImage) Display() {
    if pi.realImage == nil {
       pi.realImage = NewRealImage(pi.filename)
    }
    pi.realImage.Display()
}

func main() {
    image := NewProxyImage("test.jpg")
    // 第一次调用，会加载图片
    image.Display()
    fmt.Println()
    // 第二次调用，不会再次加载图片
    image.Display()
}
```
代码解释：

Image 是图片接口，定义了 Display 方法。
RealImage 是真实的图片对象，创建时会从磁盘加载图片。
ProxyImage 是虚拟代理，它持有图片文件名和 RealImage 的指针。在 Display 方法中，当 realImage 为 nil 时才会创建 RealImage 对象，从而实现延迟加载。
在 main 函数中，第一次调用 Display 方法会加载图片，第二次调用则不会再次加载。


保护代理
保护代理用于控制对目标对象的访问权限。以下是一个简单的保护代理示例，模拟对敏感信息的访问控制：
```go
package main

import (
    "fmt"
)

// SensitiveInfo 定义敏感信息接口
type SensitiveInfo interface {
    GetInfo() string
}

// RealSensitiveInfo 实现敏感信息接口，代表真实的敏感信息对象
type RealSensitiveInfo struct{}

func (rsi *RealSensitiveInfo) GetInfo() string {
    return "This is sensitive information."
}

// ProtectionProxy 保护代理
type ProtectionProxy struct {
    realInfo    *RealSensitiveInfo
    isAuthorized bool
}

func NewProtectionProxy() *ProtectionProxy {
    return &ProtectionProxy{
       realInfo:    &RealSensitiveInfo{},
       isAuthorized: false,
    }
}

func (pp *ProtectionProxy) Authorize() {
    pp.isAuthorized = true
    fmt.Println("User is authorized.")
}

func (pp *ProtectionProxy) GetInfo() string {
    if pp.isAuthorized {
       return pp.realInfo.GetInfo()
    }
    return "Access denied. Please authorize first."
}

func main() {
    proxy := NewProtectionProxy()
    // 未授权访问
    fmt.Println(proxy.GetInfo())
    // 授权
    proxy.Authorize()
    // 授权后访问
    fmt.Println(proxy.GetInfo())
}
```
代码解释：

SensitiveInfo 是敏感信息接口，定义了 GetInfo 方法。
RealSensitiveInfo 是真实的敏感信息对象，提供敏感信息。
ProtectionProxy 是保护代理，它持有 RealSensitiveInfo 的实例和一个授权标志。在 GetInfo 方法中，只有当用户被授权时才会返回敏感信息，否则返回访问拒绝信息。
在 main 函数中，演示了未授权和授权后访问敏感信息的不同结果。