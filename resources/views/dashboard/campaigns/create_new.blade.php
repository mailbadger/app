@extends('dashboard.layout')
@section('scripts')
    @parent
    <script type="text/javascript" src="{{asset('js/campaigns.bundle.js')}}"></script>
@endsection
@section('main')
    <h1 class="page-header">Create new campaign</h1>
    <div class="row">
        <form id="campaign-form">
            <div class="col-lg-4">
                <div class="form-group">
                    <label for="name">Campaign name:</label>
                    <input type="text" class="form-control" name="name" id="name" placeholder="Name">
                </div>
                <div class="form-group">
                    <label for="subject">Subject:</label>
                    <input type="text" class="form-control" name="subject" id="subject" placeholder="Subject">
                </div>
                <div class="form-group">
                    <label for="from-name">From name:</label>
                    <input type="text" class="form-control" name="from_name" id="from-name" placeholder="John Doe">
                </div>
                <div class="form-group">
                    <label for="from-email">From email:</label>
                    <input type="email" class="form-control" name="from_email" id="from-email" placeholder="example@foobar.com">
                </div>
            </div>
            <div class="col-lg-6 template">
                <label for="select-template">Select template:</label>
                <select id="select-template" style="width:75%">
                    <option></option>
                    @if(isset($templates) && !$templates->isEmpty())
                        @foreach($templates as $t)
                            <option value="{{$t->id}}">{{$t->name}}</option>
                        @endforeach
                    @endif
                </select>
            </div>
        </form>
    </div>
@endsection