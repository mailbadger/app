@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/components/templates-table.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">All templates</h1>
    <div class="row">
        <div class="col-lg-4">
            <a href="{{url('dashboard/new-template')}}" class="btn btn-success btn-lg">
                <span class="glyphicon glyphicon-plus"></span> Create new template
            </a>
        </div>
    </div>

    <div class="row">
        <div class="col-lg-12" id="templates"></div>
    </div>
@endsection