@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/create-new-list.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Create new subscribers list</h1>
    <div class="row" id="new-sub-list"></div>
@endsection