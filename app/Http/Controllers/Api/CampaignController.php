<?php

namespace newsletters\Http\Controllers\Api;

use Exception;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use newsletters\Http\Controllers\Controller;
use newsletters\Http\Requests;
use newsletters\Http\Requests\StoreCampaignRequest;
use newsletters\Repositories\CampaignRepository;

class CampaignController extends Controller
{
    /**
     * @var CampaignRepository
     */
    private $repository;

    public function __construct(CampaignRepository $repository)
    {
        $this->middleware('auth');

        $this->repository = $repository;
    }

    /**
     * Display a listing of the resource.
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        if ($request->has('paginate')) {
            $perPage = ($request->has('per_page')) ? $request->input('per_page') : 15;
            $campaigns = $this->repository->paginate($perPage);
        } else {
            $campaigns = $this->repository->all();
        }

        return response()->json($campaigns, 200);
    }

    /**
     * Store a newly created resource in storage.
     *
     * @param StoreCampaignRequest $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function store(StoreCampaignRequest $request)
    {
        try {
            $campaign = $this->repository->create($request->all());
            if (isset($campaign)) {
                return response()->json(['status' => 200, 'campaign' => $campaign], 200);
            }

            return response()->json(['status' => 412, 'campaign' => ['The specified resource could not be created.']],
                412);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return response()->json(['status' => 412, 'campaign' => ['The specified resource could not be created.']],
                412);
        }
    }

    /**
     * Display the specified resource.
     *
     * @param  int $id
     * @return \Illuminate\Http\JsonResponse
     */
    public function show($id)
    {
        try {
            $campaign = $this->repository->find($id);

            return response()->json($campaign, 200);
        } catch (ModelNotFoundException $e) {
            Log::error($e->getMessage());

            return response()->json(['status' => 404, 'message' => 'The specified resource does not exist.'], 404);
        }
    }

    /**
     * Update the specified resource in storage.
     *
     * @param  Request $request
     * @param  int $id
     * @return Response
     */
    public function update(Request $request, $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     *
     * @param  int $id
     * @return Response
     */
    public function destroy($id)
    {
        try {
            $this->repository->delete($id);

            return response()->json(['status' => 200, 'message' => 'The specified resource has been deleted.'], 200);
        } catch (Exception $e) {
            Log::error($e->getMessage());

            return response()->json(['status' => 422, 'campaign' => ['The specified resource could not be deleted.']],
                422);
        }
    }
}
