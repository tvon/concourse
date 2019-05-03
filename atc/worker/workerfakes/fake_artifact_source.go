// Code generated by counterfeiter. DO NOT EDIT.
package workerfakes

import (
	"io"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/concourse/atc/worker"
)

type FakeArtifactSource struct {
	StreamFileStub        func(lager.Logger, string) (io.ReadCloser, error)
	streamFileMutex       sync.RWMutex
	streamFileArgsForCall []struct {
		arg1 lager.Logger
		arg2 string
	}
	streamFileReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	streamFileReturnsOnCall map[int]struct {
		result1 io.ReadCloser
		result2 error
	}
	StreamToStub        func(lager.Logger, worker.ArtifactDestination) error
	streamToMutex       sync.RWMutex
	streamToArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.ArtifactDestination
	}
	streamToReturns struct {
		result1 error
	}
	streamToReturnsOnCall map[int]struct {
		result1 error
	}
	VolumeOnStub        func(lager.Logger, worker.Worker) (worker.Artifact, bool, error)
	volumeOnMutex       sync.RWMutex
	volumeOnArgsForCall []struct {
		arg1 lager.Logger
		arg2 worker.Worker
	}
	volumeOnReturns struct {
		result1 worker.Artifact
		result2 bool
		result3 error
	}
	volumeOnReturnsOnCall map[int]struct {
		result1 worker.Artifact
		result2 bool
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArtifactSource) StreamFile(arg1 lager.Logger, arg2 string) (io.ReadCloser, error) {
	fake.streamFileMutex.Lock()
	ret, specificReturn := fake.streamFileReturnsOnCall[len(fake.streamFileArgsForCall)]
	fake.streamFileArgsForCall = append(fake.streamFileArgsForCall, struct {
		arg1 lager.Logger
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("StreamFile", []interface{}{arg1, arg2})
	fake.streamFileMutex.Unlock()
	if fake.StreamFileStub != nil {
		return fake.StreamFileStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.streamFileReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeArtifactSource) StreamFileCallCount() int {
	fake.streamFileMutex.RLock()
	defer fake.streamFileMutex.RUnlock()
	return len(fake.streamFileArgsForCall)
}

func (fake *FakeArtifactSource) StreamFileCalls(stub func(lager.Logger, string) (io.ReadCloser, error)) {
	fake.streamFileMutex.Lock()
	defer fake.streamFileMutex.Unlock()
	fake.StreamFileStub = stub
}

func (fake *FakeArtifactSource) StreamFileArgsForCall(i int) (lager.Logger, string) {
	fake.streamFileMutex.RLock()
	defer fake.streamFileMutex.RUnlock()
	argsForCall := fake.streamFileArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeArtifactSource) StreamFileReturns(result1 io.ReadCloser, result2 error) {
	fake.streamFileMutex.Lock()
	defer fake.streamFileMutex.Unlock()
	fake.StreamFileStub = nil
	fake.streamFileReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifactSource) StreamFileReturnsOnCall(i int, result1 io.ReadCloser, result2 error) {
	fake.streamFileMutex.Lock()
	defer fake.streamFileMutex.Unlock()
	fake.StreamFileStub = nil
	if fake.streamFileReturnsOnCall == nil {
		fake.streamFileReturnsOnCall = make(map[int]struct {
			result1 io.ReadCloser
			result2 error
		})
	}
	fake.streamFileReturnsOnCall[i] = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeArtifactSource) StreamTo(arg1 lager.Logger, arg2 worker.ArtifactDestination) error {
	fake.streamToMutex.Lock()
	ret, specificReturn := fake.streamToReturnsOnCall[len(fake.streamToArgsForCall)]
	fake.streamToArgsForCall = append(fake.streamToArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.ArtifactDestination
	}{arg1, arg2})
	fake.recordInvocation("StreamTo", []interface{}{arg1, arg2})
	fake.streamToMutex.Unlock()
	if fake.StreamToStub != nil {
		return fake.StreamToStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.streamToReturns
	return fakeReturns.result1
}

func (fake *FakeArtifactSource) StreamToCallCount() int {
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	return len(fake.streamToArgsForCall)
}

func (fake *FakeArtifactSource) StreamToCalls(stub func(lager.Logger, worker.ArtifactDestination) error) {
	fake.streamToMutex.Lock()
	defer fake.streamToMutex.Unlock()
	fake.StreamToStub = stub
}

func (fake *FakeArtifactSource) StreamToArgsForCall(i int) (lager.Logger, worker.ArtifactDestination) {
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	argsForCall := fake.streamToArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeArtifactSource) StreamToReturns(result1 error) {
	fake.streamToMutex.Lock()
	defer fake.streamToMutex.Unlock()
	fake.StreamToStub = nil
	fake.streamToReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactSource) StreamToReturnsOnCall(i int, result1 error) {
	fake.streamToMutex.Lock()
	defer fake.streamToMutex.Unlock()
	fake.StreamToStub = nil
	if fake.streamToReturnsOnCall == nil {
		fake.streamToReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.streamToReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeArtifactSource) VolumeOn(arg1 lager.Logger, arg2 worker.Worker) (worker.Artifact, bool, error) {
	fake.volumeOnMutex.Lock()
	ret, specificReturn := fake.volumeOnReturnsOnCall[len(fake.volumeOnArgsForCall)]
	fake.volumeOnArgsForCall = append(fake.volumeOnArgsForCall, struct {
		arg1 lager.Logger
		arg2 worker.Worker
	}{arg1, arg2})
	fake.recordInvocation("VolumeOn", []interface{}{arg1, arg2})
	fake.volumeOnMutex.Unlock()
	if fake.VolumeOnStub != nil {
		return fake.VolumeOnStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.volumeOnReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeArtifactSource) VolumeOnCallCount() int {
	fake.volumeOnMutex.RLock()
	defer fake.volumeOnMutex.RUnlock()
	return len(fake.volumeOnArgsForCall)
}

func (fake *FakeArtifactSource) VolumeOnCalls(stub func(lager.Logger, worker.Worker) (worker.Artifact, bool, error)) {
	fake.volumeOnMutex.Lock()
	defer fake.volumeOnMutex.Unlock()
	fake.VolumeOnStub = stub
}

func (fake *FakeArtifactSource) VolumeOnArgsForCall(i int) (lager.Logger, worker.Worker) {
	fake.volumeOnMutex.RLock()
	defer fake.volumeOnMutex.RUnlock()
	argsForCall := fake.volumeOnArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeArtifactSource) VolumeOnReturns(result1 worker.Artifact, result2 bool, result3 error) {
	fake.volumeOnMutex.Lock()
	defer fake.volumeOnMutex.Unlock()
	fake.VolumeOnStub = nil
	fake.volumeOnReturns = struct {
		result1 worker.Artifact
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeArtifactSource) VolumeOnReturnsOnCall(i int, result1 worker.Artifact, result2 bool, result3 error) {
	fake.volumeOnMutex.Lock()
	defer fake.volumeOnMutex.Unlock()
	fake.VolumeOnStub = nil
	if fake.volumeOnReturnsOnCall == nil {
		fake.volumeOnReturnsOnCall = make(map[int]struct {
			result1 worker.Artifact
			result2 bool
			result3 error
		})
	}
	fake.volumeOnReturnsOnCall[i] = struct {
		result1 worker.Artifact
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeArtifactSource) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.streamFileMutex.RLock()
	defer fake.streamFileMutex.RUnlock()
	fake.streamToMutex.RLock()
	defer fake.streamToMutex.RUnlock()
	fake.volumeOnMutex.RLock()
	defer fake.volumeOnMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeArtifactSource) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ worker.ArtifactSource = new(FakeArtifactSource)
