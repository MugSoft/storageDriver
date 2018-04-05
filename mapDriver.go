package storageDriver

import (
	"fmt"
	"io"
	"reflect"
	"sync"
)

type mapDriver struct {
	database   string
	collection string
	sync.Mutex
	store map[string]map[string][]Document
}

func (m *mapDriver) Driver() (StorageDriver, error) {
	if d.database == "" || d.collection == "" {
		return nil, fmt.Errorf("database or collection is not set")
	}
	return d, nil
}
func (m *mapDriver) DB(name string) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}
	m.database = name
	return nil
}
func (m *mapDriver) Table(name string) error {

	if name == "" {
		return fmt.Errorf("empty name")
	}
	m.collection = name
	return nil
}

func (m *mapDriver) Clone() Meta {
	var cpy = *m
	return &cpy
}
func (d *mapDriver) Get(Query Document) ([]Document, error) {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	var docs = make([]Document, 0)
	var match bool = true
	d.Lock()
	defer d.Unlock()
	for _, DBDoc := range d.store[d.database][d.collection] {
		for k, v := range Query {
			if val, ok := DBDoc[k]; !ok || !reflect.DeepEqual(val, v) {
				match = false
				goto next
			}
		}
		if match {
			docs = append(docs, DBDoc)
		}
	next:
		match = true
	}
	var err error
	if len(docs) == 0 {
		err = fmt.Errorf("no documents found")
	}
	return docs, err
}

func (d *mapDriver) Insert(doc Document) error {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	d.Lock()
	defer d.Unlock()
	d.store[d.database][d.collection] = append(d.store[d.database][d.collection], doc)
	return nil
}

func (d *mapDriver) GetOne(Query Document) (Document, error) {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	var match bool = true
	d.Lock()
	defer d.Unlock()
	for _, DBDoc := range d.store[d.database][d.collection] {
		for k, v := range Query {
			if val, ok := DBDoc[k]; !ok || !reflect.DeepEqual(val, v) {
				match = false
				goto next
			}
		}
		if match {
			return DBDoc, nil
		}
	next:
		match = true
	}
	return nil, fmt.Errorf("no documents found")
}
func (d *mapDriver) Custom(_ interface{}) ([]Document, error) {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	return nil, fmt.Errorf("not implemented yet")
}

func (d *mapDriver) InsertMulti(docs []Document) error {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	for _, doc := range docs {
		d.Insert(doc)
	}
	return nil
}

func (d *mapDriver) InsertMultiNoFail(docs []Document, _ ...io.Writer) []error {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	d.InsertMulti(docs)
	return nil
}

func (d *mapDriver) Update(Query Document, UpdatedFields Document) error {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	doc, err := d.GetOne(Query)
	if nil != err {
		return err
	}
	d.Lock()
	defer d.Unlock()
	for k, v := range UpdatedFields {
		doc[k] = v
	}
	return nil
}
func (d *mapDriver) UpdateMulti(Query, UpdatedFields Document) (int, error) {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	docs, err := d.Get(Query)
	if nil != err {
		return 0, err
	}
	d.Lock()
	defer d.Unlock()
	for _, doc := range docs {
		for k, v := range UpdatedFields {
			doc[k] = v
		}
	}
	return len(docs), nil
}
func (d *mapDriver) Save(Query, Doc Document) error {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	doc, err := d.GetOne(Query)
	if nil != err {
		dd := make(Document)
		for k, v := range Query {
			dd[k] = v
		}
		for k, v := range Doc {
			dd[k] = v
		}
		return d.Insert(dd)
	}
	return d.Update(doc, Doc)

	return nil
}
func (d *mapDriver) Remove(Query Document) error {
	if _, ok := d.store[d.database]; !ok {
		d.store[d.database] = make(map[string][]Document)
	}
	var match bool = true
	d.Lock()
	defer d.Unlock()
	for docIndex, DBDoc := range d.store[d.database][d.collection] {
		for k, v := range Query {
			if val, ok := DBDoc[k]; !ok || !reflect.DeepEqual(val, v) {
				match = false
				goto next
			}
		}
		if match {
			d.store[d.database][d.collection] = append(d.store[d.database][d.collection][:docIndex], d.store[d.database][d.collection][docIndex+1:]...)
			return nil
		}
	next:
		match = true
	}

	return fmt.Errorf("no document removed")
}