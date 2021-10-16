### 使用 io 包下的 Copy 函数

参考代码：./v3/copy_file.go

**关键**使用的**代码**就是：

~~~go
...
for num += 1; num < endNum; num++ {
    name := strconv.Itoa(num) + suffix
    // create a new file
    dstFile, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        log.Fatal(err)
        return
    }
    defer dstFile.Close()

    io.Copy(dstFile, file)
    // must be call seek to set the offset for next read!
    file.Seek(0, io.SeekStart)
}
...
~~~

`file.Seek(0, io.SeekStart)` 这行代码是**必需的**，否则，在 for 循环的后续无法从 file（也就是 src）中读取到源数据。file 读取的位置游标随着 read 操作依次向 EOF 方向移动，下一次 Read 操作会从游标位置开始读取。

另外，`defer dstFile.Close()` 在 for 循环中被调用时，会在本函数执行完成后依次添加 defer 语句。而且，for 中最先执行的 defer 语句内容的执行时机会是最后，是一个**倒置**的逻辑。

标准库 io.Copy 函数：

~~~go
// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}
~~~



