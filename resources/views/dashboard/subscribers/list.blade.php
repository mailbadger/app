@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/subscribers.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Lists</h1>
    <div id="sub-lists"></div>
@endsection
