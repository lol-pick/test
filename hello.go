package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Rope struct {
	left  *Rope
	right *Rope
	len   int
	start int // 0-индексный старт в исходной строке s
	end   int // 0-индексный конец в исходной строке s
}

func (r *Rope) isLeaf() bool {
	return r.left == nil && r.right == nil
}

// Получение символа на позиции pos (1-индексно)
func (r *Rope) getChar(pos int) byte {
	cur := r
	p := pos
	for !cur.isLeaf() {
		if p <= cur.left.len {
			cur = cur.left
		} else {

			cur = cur.right
		}
	}
	// p — позиция внутри листа (1-индексно)
	// Преобразуем в 0-индексный индекс в исходной строке
	idx := cur.start + (p - 1)
	return s[idx]
}

// Разделение верёвки на две части: левая содержит первые `pos` символов
func (r *Rope) split(pos int) (*Rope, *Rope) {
	if pos <= 0 {
		return nil, r
	}
	if pos >= r.len {
		return r, nil
	}
	if r.isLeaf() {
		// Разделяем лист на два новых листа
		mid := r.start + pos - 1 // конец левой части (0-индексно)
		left := &Rope{
			start: r.start,
			end:   mid,
			len:   pos,
		}
		right := &Rope{
			start: mid + 1,
			end:   r.end,
			len:   r.len - pos,
		}
		return left, right
	}
	if pos <= r.left.len {
		left1, right1 := r.left.split(pos)
		rightPart := concat(right1, r.right)
		return left1, rightPart
	} else {
		left2, right2 := r.right.split(pos - r.left.len)
		leftPart := concat(r.left, left2)
		return leftPart, right2
	}
}

func concat(a, b *Rope) *Rope {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return &Rope{
		left:  a,
		right: b,
		len:   a.len + b.len,
	}
}

var s string

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Чтение n и q
	line, _ := reader.ReadString('\n')
	parts := strings.Fields(line)
	n, _ := strconv.Atoi(parts[0])
	q, _ := strconv.Atoi(parts[1])

	// Чтение строки
	s, _ = reader.ReadString('\n')
	s = strings.TrimSpace(s)

	// Инициализация верёвки
	root := &Rope{
		start: 0,
		end:   n - 1,
		len:   n,
	}

	var output []byte
	for i := 0; i < q; i++ {
		line, _ = reader.ReadString('\n')
		parts := strings.Fields(line)

		if parts[0] == "1" {
			// Операция 1 l r: скопировать подстроку [l, r] и добавить в конец
			l, _ := strconv.Atoi(parts[1])
			r, _ := strconv.Atoi(parts[2])

			// Извлекаем подстроку [l, r]
			left, temp := root.split(l - 1)     // left = [1, l-1]
			mid, right := temp.split(r - l + 1) // mid = [l, r], right = [r+1, end]

			// Восстанавливаем исходную строку
			original := concat(concat(left, mid), right)

			// Добавляем копию подстроки в конец
			root = concat(original, mid)
		} else {
			// Операция 2 i: запрос символа на позиции i
			idx, _ := strconv.Atoi(parts[1])
			c := root.getChar(idx)
			output = append(output, c)
			output = append(output, ' ')
		}
	}

	// Удаляем последний пробел
	if len(output) > 0 {
		output = output[:len(output)-1]
	}
	os.Stdout.Write(output)
}
