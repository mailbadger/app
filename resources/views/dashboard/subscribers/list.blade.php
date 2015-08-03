@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/sub-list.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">All lists</h1>
    <div id="sub-lists"></div>
@endsection