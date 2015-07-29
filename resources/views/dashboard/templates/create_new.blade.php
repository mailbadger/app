@extends('dashboard.layout')
@section('scripts')
    @parent
    <script src="//cdn.ckeditor.com/4.5.1/full/ckeditor.js"></script>
    <script type="text/javascript" src="{{asset('js/components/templates/templates-form.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Create new template</h1>
    <div class="row" id="new-template"></div>
@endsection